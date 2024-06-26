#%%

import os,sys,math

class Adb:
    BACK=4
    ENTER=66
    HOME=3
    POWER=26

    def __init__(self, adb_path):
        self.adb_path = adb_path

    def screencast(self):
        os.system(f'{self.adb_path} shell screencap /sdcard/screen.png')
        os.system(f'{self.adb_path} pull /sdcard/screen.png static/screen.jpg')
        from PIL import Image
        img_png = Image.open('static/screen.jpg').convert('RGB')
        img_png.save('static/screen.jpg')

    def tap(self, x, y):
        return os.system(f'{self.adb_path} shell input tap {x} {y}')

    def swipe(self, x1,y1,x2,y2,dt=100):
        return os.system(f'{self.adb_path} shell input swipe {x1} {y1} {x2} {y2} {dt}')

    def key(self, k):
        return os.system(f'{self.adb_path} shell input keyevent {k}')
    
    def run(self, app):
        return os.system(f'{self.adb_path} shell am start {app}')

