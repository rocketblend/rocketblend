import bpy
import sys
import ast
import os
import argparse
import logging
import addon_utils

from pathlib import Path
from logging.handlers import RotatingFileHandler

if not os.path.exists("logs"):
    os.makedirs("logs")

logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s %(message)s',
    datefmt='%a, %d %b %Y %H:%M:%S',
    handlers=[
        # RotatingFileHandler(
        #     "logs/rocketblend.log",
        #     mode="a",
        #     maxBytes=5*1024*1024,
        #     backupCount=2,
        #     encoding=None,
        #     delay=0
        # ),
        logging.StreamHandler()
]) 

log = logging.getLogger(__name__)

# This is the script that is run when Blender is started from the CLI.
# It is used to load the addons specified in the command line.
# The script is called with the following arguments:
# -a, --addons: a comma-separated list of addons to load.
#       E.g., -a "[{"name": "addon1", version: "0.1.0", "path": "C:\Users\user\Documents\blender\addons\addon1"},
#                  {"name": "addon2", "path": "C:\Users\user\Documents\blender\addons\addon2"}]"
#  If the addon is already installed on this build, it will be
#  loaded. Otherwise, it will be installed and then loaded.

class ArgumentParserForBlender(argparse.ArgumentParser):
    """
    This class is identical to its superclass, except for the parse_args
    method (see docstring). It resolves the ambiguity generated when calling
    Blender from the CLI with a python script, and both Blender and the script
    have arguments. E.g., the following call will make Blender crash because
    it will try to process the script's -a and -b flags:
    >>> blender --python my_script.py -a 1 -b 2

    To bypass this issue this class uses the fact that Blender will ignore all
    arguments given after a double-dash ('--'). The approach is that all
    arguments before '--' go to Blender, arguments after go to the script.
    The following calls work fine:
    >>> blender --python my_script.py -- -a 1 -b 2
    >>> blender --python my_script.py --
    """

    def _get_argv_after_doubledash(self):
        """
        Given the sys.argv as a list of strings, this method returns the
        sublist right after the '--' element (if present, otherwise returns
        an empty list).
        """
        try:
            idx = sys.argv.index("--")
            return sys.argv[idx+1:] # the list after '--'
        except ValueError as e: # '--' not in the list:
            return []

    # overrides superclass
    def parse_args(self):
        """
        This method is expected to behave identically as in the superclass,
        except that the sys.argv list will be pre-processed using
        _get_argv_after_doubledash before. See the docstring of the class for
        usage examples and details.
        """
        return super().parse_args(args=self._get_argv_after_doubledash())

class Addon(object):
    """
    This class represents an addon, and is used to parse the command line
    argument.
    """
    def __init__(self, name: str, version: str, path: str):
        self.name = name
        self.version = self._parse_version(version)
        self.path = path # Empty paths are used to indicate that the addon is pre-installed.

    def __str__(self):
        return f"Addon(name={self.name}, version={self.version}, path={self.path})"

    def __repr__(self):
        return self.__str__()

    def _parse_version(self, version_str: str) -> tuple[int, int, int]:
        if version_str == "":
            return (-1, -1, -1)
        
        version = ()
        for element in version_str.split('.'):
            version += (int(element),)
        
        return version

class AddonManager(object):
    def __init__(self, addons: list[dict]):
        self.addons = []
        for addon in addons:
            path = Path(addon.get("path", ""))
            if path.exists() or path == "":
                name = addon.get("name", "")
                version = addon.get("version", "")
                self.addons.append(Addon(str(name), str(version), addon.get("path", "")))

    def get(self, ignore_pre_installed: bool = False) -> list[Addon]:
        return [addon for addon in self.addons if addon.path != "" or not ignore_pre_installed]

    def find(self, name: str) -> Addon:
        for addon in self.addons:
            if addon.name == name:
                return addon
        return None

class Startup():
    """
    This class is called at Blender startup, and is used to install/enable the
    addons specified by the command line argument.
    """
    
    def __init__(self, addons: list[Addon]):
        logging.debug(f"Starting Blender with the following addons: {addons}")

        self.manager = AddonManager(addons)

        self.install_addons(True)
        self.reset_addons()

        logging.debug(f"Finished loading addons")

    def install_addons(self, overwrite: bool = False) -> None:
        """
        Installs any addon not already installed on this build.
        """
        installed = []

        for mod in addon_utils.modules():
            addon = self.manager.find(mod.__name__)

            if addon is not None and addon.version == mod.bl_info.get('version', (0, 0, 0)):
                installed.append(mod.__name__)

        for addon in self.manager.get(ignore_pre_installed=True):
            if addon.name not in installed:
                logging.debug(f"Installing addon {addon.name} {addon.version} from {addon.path}")
                bpy.ops.preferences.addon_install(filepath=addon.path, overwrite=overwrite)

    def reset_addons(self) -> None:
        """
        Resets addons to only the ones defined.
        """  
        enable = [addon.name for addon in self.manager.get()]

        for addon in bpy.context.preferences.addons:
            if addon.module not in enable:
                self.disable_addon(addon.module)
            else:
                enable.remove(addon.module)

        for addonName in enable:
            self.enable_addon(addonName)

    def enable_addon(self, addon_name: str) -> None:
        """
        Enables the addon with the given name.
        """
        mod = addon_utils.enable(addon_name, default_set=False, persistent=False)

        if mod:
            logging.debug(f"Enabled addon {addon_name}")

            info = addon_utils.module_bl_info(mod)
            info_ver = info.get("blender", (0, 0, 0))

            if info_ver > bpy.app.version:
                logging.debug(f"Addon {addon_name} was written in Blender {info_ver[0]}.{info_ver[1]} and may not work correctly in Blender {bpy.app.version[0]}.{bpy.app.version[1]}")
        else:
            logging.debug(f"Failed to enable addon {addon_name}")

    def disable_addon(self, addon_name: str) -> None:
        """
        Disables the addon with the given name.
        """
        addon_utils.disable(addon_name, default_set=False)
        logging.debug(f"Disabled addon {addon_name}")

parser = ArgumentParserForBlender()
parser.add_argument("-a", "--addons", help="Addons to load", type=ast.literal_eval, default={})

args = parser.parse_args()

Startup(args.addons)