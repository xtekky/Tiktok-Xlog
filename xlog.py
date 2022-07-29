#!usr/bin/env python
#-*- coding:utf-8 -*-

"""
@price: 500€
@author: t.me/xtekky
@file: xlog.py
@state: preview
"""

import binascii
import codecs
import ctypes

class Xlog:
    def encode(self, data):
        # removed from preview

        return encoded

    def decode(self, data: hex):
        # removed from preview

        return decoded

    def __calculate(self, input, times):
        # removed from preview

    def __xor(self, x, y):
        # removed from preview
        return format(
          (
            int(
              response
              , 2
            )
          )
          , 'x'
        )

    def __handle(self, hex):
        # removed from preview

    def __shift(self, point):
       # removed from preview

    def  __calcrev(self, data, num):
      s108 = 0xbfffe920 << 0 & 0xFFFFFFFF
      s136 = 0x0
      s140 = int(data[0:8], 16) << 0 & 0xFFFFFFFF
      s144 = int(data[8:16], 16) << 0 & 0xFFFFFFFF

      for i in range(num):
          r2 = s108
          r6 = s136
          r4 = s144
          r5 = r6 & 3 & 0xFFFFFFFF
          r0 = r4 << 4 & 0xFFFFFFFF
          r2 = self.getShifting(r2 + (r5 << 2) & 0xFFFFFFFF)
          r0 = ((r0 ^ (self.rshift(r4, 5))) + r4) << 0
          r2 = ctypes.c_int((r2 + r6) << 0 ^ 0).value
          r0 = r0 ^ r2
          s140 = ctypes.c_int((s140 + r0) << 0 ^ 0).value
          s136 = ctypes.c_int((s136 - 0x61c88647) << 0 ^ 0).value

          r5 = s108
          r4 = s140
          r2 = s140
          r0 = s140
          r6 = s136
          s = format(self.rshift((r6 >> 0xb), 0), 'b')
          if len(s) < 3:
              s = "0"
          else:
              s = s[len(s) - 2:]

          r6 = int(s, 2)
          r0 = ctypes.c_int(((self.rshift(r2, 5) ^ r0 << 4) + r4) << 0).value
          r5 = self.getShifting(r5 + (r6 << 2))
          r2 = ctypes.c_int((s136 + r5) << 0 ^ 0).value
          r0 = r0 ^ r2
          s144 = ctypes.c_int((s144 + r0) << 0 ^ 0).value

      str140 = format(self.rshift(s140, 0), 'x')

      str144 = format(self.rshift(s144, 0), 'x')


      if len(str140) < 8:
          count = 8 - len(str140)
          for i in range(count):
              str140 = "0" + str140

      if len(str144) < 8:
          count = 8 - len(str144)
          for i in range(count):
              str144 = "0" + str144

      return str140 + str144

    def __reverse(self, hex: str):
        return hex[6:8] + hex[4:6] + hex[2:4] + hex[0:2]

    def __rshift(self, val, n):
        return (val % 0x100000000) >> n

    def __hxstr(self, num: int):
        s = format(num, 'x')
        if len(s) < 2:
            return '0' + s
        return s

    def __fch(self, xlog):
        xlog = xlog[0:len(xlog) - 21]
        fch_str = binascii.crc32(xlog.encode("utf-8"))
        fch_str = str(fch_str)

        for i in range(len(fch_str), 10):
            fch_str = '0' + fch_str

        return fch_str
