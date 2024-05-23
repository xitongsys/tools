#%%

import time
from datetime import datetime

from adb import Adb

if __name__ == '__main__':
    now = datetime.now()
    if now.hour not in [7,8,18,19]:
        exit(0)

    phone = Adb(r'C:\Users\zxt\AppData\Local\Android\Sdk\platform-tools\adb.exe')

    for i in range(4):
        phone.key(Adb.POWER)
        #phone.swipe(0,0,1000,1000)
        phone.run('com.ss.android.lark')
        time.sleep(20)
