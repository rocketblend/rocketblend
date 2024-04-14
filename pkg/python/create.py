import bpy
bpy.ops.wm.save_as_mainfile(filepath=r'{{ .filePath }}')
bpy.ops.wm.quit_blender()