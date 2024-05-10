#%%

import os,sys,math
import json
from flask import Flask, request, render_template,jsonify
import cv2
from threading import Timer

class Adb:
    BACK=4
    ENTER=66
    HOME=3
    POWER=26

    def __init__(self, adb_path):
        self.adb_path = adb_path

    def screencast(self):
        os.system(f'{self.adb_path} shell screencap /sdcard/screen.png')
        return os.system(f'{self.adb_path} pull /sdcard/screen.png static/screen.png')

    def tap(self, x, y):
        return os.system(f'{self.adb_path} shell input tap {x} {y}')

    def swipe(self, x1,y1,x2,y2):
        return os.system(f'{self.adb_path} shell input swipe {x1} {y1} {x2} {y2}')

    def key(self, k):
        return os.system(f'{self.adb_path} shell input keyevent {k}')


class WebServer:
    def __init__(self, ip, port):
        self.ip = ip
        self.port = port
        self.adb = Adb(r'C:\Users\zxt\AppData\Local\Android\Sdk\platform-tools\adb.exe')
        self.adb.screencast()
        self.h, self.w, _ = cv2.imread('static/screen.png').shape
            
    def run(self):
        app = Flask(__name__)
        app.config["TEMPLATES_AUTO_RELOAD"] = True

        def refresh_timer():
            self.adb.screencast()
            timer = Timer(0.5, refresh_timer)
            timer.start()
        refresh_timer()     

        @app.route('/')
        def index():
            return render_template('index.html')
        

        @app.route('/refresh',methods=['GET'])
        def refresh():
            return jsonify(self.adb.screencast())
        

        @app.route("/tap",methods=['GET'])
        def tap():
            x,y = int(float(request.args.get('x'))*self.w), int(float(request.args.get('y'))*self.h)
            res = self.adb.tap(x, y)
            return jsonify(res)
        
        @app.route("/swipe",methods=['GET'])
        def swipe():
            x1,y1 = int(float(request.args.get('x1'))*self.w), int(float(request.args.get('y1'))*self.h)
            x2,y2 = int(float(request.args.get('x2'))*self.w), int(float(request.args.get('y2'))*self.h)
            res = self.adb.swipe(x1, y1, x2, y2)
            return jsonify(res)


        app.run(self.ip, self.port)





if __name__ == '__main__':
    webserver = WebServer('0.0.0.0',4001)
    webserver.run()
