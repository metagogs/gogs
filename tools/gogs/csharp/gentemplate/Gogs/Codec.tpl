using System;
using System.Collections.Generic;
using Google.Protobuf;

namespace Gogs
{
    class ProtoEncode : IEncode
    {
        public byte[] Encode(IMessage msg)
        {
            if (msg is IMessage message)
            {
                return message.ToByteArray();
            }

            throw new Exception("object is not the message type");
        }
    }

    class ProtoDecode : IDecode
    {
        static readonly Dictionary<String, MessageParser> parses = new Dictionary<String, MessageParser>();

        public T Decode<T>(byte[] data) where T : IMessage<T>, new()
        {
            if (parses.TryGetValue(typeof(T).Name, out MessageParser parse))
            {
                return (T)parse.ParseFrom(data);
            }

            Google.Protobuf.MessageParser<T> newParse = new Google.Protobuf.MessageParser<T>(() => new T());
            parses[typeof(T).Name] = newParse;
            return newParse.ParseFrom(data);
        }
    }
}