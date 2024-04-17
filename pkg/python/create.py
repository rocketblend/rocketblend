import bpy
bpy.ops.wm.save_as_mainfile(filepath=r'{{ .FilePath }}')
bpy.ops.wm.quit_blender()