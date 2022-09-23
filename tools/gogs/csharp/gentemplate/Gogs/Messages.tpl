using System;
using System.Collections.Generic;
using System.Reflection;
using Google.Protobuf;

namespace Gogs
{
    class Messages
    {

        public const int CODEC_JSON_DATA_NO_HEADER = 0;
        public const int CODEC_JSON_DATA = 1;
        public const int CODEC_PROTO_DATA = 2;


        private static Messages instance = null;

        public static IEncode _encode = new ProtoEncode();
        public static int _encodeType = CODEC_PROTO_DATA;
        public static Dictionary<int, IDecode> _decodes = new Dictionary<int, IDecode>();


        public static Messages Instance()
        {
            if (instance == null)
            {
                instance = new Messages();
                instance.AddDecode(CODEC_PROTO_DATA, new ProtoDecode());
                // init common message
                Dispatch.AddField("Ping", typeof(Ping), Packet.CreateAction(1, 1, 1));
                Dispatch.AddField("Pong", typeof(Pong), Packet.CreateAction(1, 1, 2));
            }
            return instance;
        }

        public static void Message(string name, byte[] data)
        {
            if(Instance().Decode(data,out Packet packet))
            {
                EventsManager.Broadcast(name, (IMessage)packet.obj);
            }
        }

        public void SetEncode(IEncode e, int type)
        {
            _encode = e;
            _encodeType = type;

        }

        public void AddDecode(int t, IDecode d)
        {
            _decodes[t] = d;
        }

        public bool Encode(Object msg, out Packet packet)
        {
            packet = null;
            int action = Dispatch.GetActionByName(msg.GetType().Name);
            if(msg is IMessage message)
            {
                byte[] data = _encode.Encode(message);
                packet = Packet.NewPacketWithHeader(data, 1, _encodeType, action);
                return true;
            }
            
            return false;
            
        }

        public bool Decode(byte[] data, out Packet packet)
        {
            if (!Packet.ParsePacket(data, out packet))
            {
                return false;
            }

            if(_decodes.TryGetValue(packet.GetEncodeType(), out IDecode decode))
            {
                Type t = Dispatch.GetTypeByAction(packet.GetActionKey());

                MethodInfo mi = decode.GetType().GetMethod("Decode").MakeGenericMethod(new Type[] { t });
                packet.obj = mi.Invoke(decode, new object[]{ packet.data });
                return true;
            }

            return false;
        }

    }

    class Dispatch
    {
        static readonly Dictionary<String, int> nameWithAction = new Dictionary<String, int>();
        static readonly Dictionary<int, Type> actionWithName = new Dictionary<int, Type>();

        public static void AddField(String fieldName, Type t, int action)
        {
            nameWithAction[fieldName] = action;
            actionWithName[action] = t;
        }

        public static int GetActionByName(String fieldName)
        {
            return nameWithAction[fieldName];
        }

        public static Type GetTypeByAction(int action)
        {
            return actionWithName[action];
        }
    }
}
