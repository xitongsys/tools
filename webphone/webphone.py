#%%

import os,sys,math
import json
from flask import Flask, request, render_template,jsonify
import cv2
from threading import Timer

from adb import Adb


class WebServer:
    def __init__(self, ip, port):
        self.ip = ip
        self.port = port
        self.adb = Adb(r'C:\Users\zxt\AppData\Local\Android\Sdk\platform-tools\adb.exe')
        self.adb.screencast()
        self.h, self.w, _ = cv2.imread('static/screen.jpg').shape
            
    def run(self):
        app = Flask(__name__)
        app.config["TEMPLATES_AUTO_RELOAD"] = True

        def refresh_timer():
            self.adb.screencast()
            timer = Timer(1, refresh_timer)
            timer.start()
        refresh_timer()     


        @app.route('/')
        def index():
            return render_template('index.html')

        @app.route('/power',methods=['GET'])
        def power():
            return jsonify(self.adb.key(Adb.POWER))


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
            dt = int(request.args.get('dt'))
            res = self.adb.swipe(x1, y1, x2, y2, dt)
            return jsonify(res)


        app.run(self.ip, self.port)





if __name__ == '__main__':
    webserver = WebServer('0.0.0.0',4001)
    webserver.run()
