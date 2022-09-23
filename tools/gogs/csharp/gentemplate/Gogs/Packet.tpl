﻿using System;
using System.Collections.Generic;

namespace Gogs
{
    public class Packet
    {
        public byte[] header;
        public byte[] data;
        public Object obj;

        public const int HeaderLength = 8;
        public const int HeaderFlag = 0X7e;


        public static int CreateAction(int packetType, int module, int action)
        {
            return packetType << 22 | module << 16 | action;
        }

        public static bool ParsePacket(byte[] data, out Packet packet)
        {
            packet = null;
            if (data.Length < 8)
            {
                return false;
            }
            if(data.Length > 1 << 24)
            {
                return false;
            }
            byte[] header = data[0..8];
            if(header[0] != HeaderFlag)
            {
                return false;
            }
            packet = new Packet();
            packet.header = header;
            packet.data = data[8..data.Length];

            return true;
        }

        public static Packet NewPacketWithHeader(byte[] data, int version, int encodeType, int action)
        {
            byte[] header = new byte[8];
            header[0] = 0x7E;
            header[1] = (byte)(version << 3 | encodeType);
            header[2] = (byte)(action >> 16);
            header[3] = (byte)(action >> 8);
            header[4] = (byte)(action);

            int len = data.Length;
            header[5] = (byte)(len >> 16);
            header[6] = (byte)(len >> 8);
            header[7] = (byte)(len);



            Packet packet = new Packet();
            packet.header = header;
            packet.data = data;

            return packet;

        }

        public byte[] ToByteArray()
        {
            List<byte> packet = new List<byte>();
            packet.AddRange(header);
            packet.AddRange(data);
            return packet.ToArray();
        }


        public int GetEncodeType()
        {
            return header[1] & (0xff >> 5);
        }

        public int GetActionKey()
        {
            return (int)(header[2] << 16) | (int)(header[3] << 8) | (int)header[4];
        }
    }

}
