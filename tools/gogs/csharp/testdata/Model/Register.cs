using System;
using Gogs;
using Google.Protobuf;
namespace Model
{
    public static class GInit
    {
        public const int SERVICE_TYPE = 2;

        public static void Init()
        {
            
            Dispatch.AddField("BindUser", typeof(BindUser), Packet.CreateAction(SERVICE_TYPE, 1, 1));  // 0x810001 8454145 
            
            Dispatch.AddField("BindSuccess", typeof(BindSuccess), Packet.CreateAction(SERVICE_TYPE, 1, 2));  // 0x810002 8454146 
            
            Gogs.Messages.Instance();
        }
    }

    public static class GMessages
    {
        public static void Message(string name, byte[] data)
        {
            Gogs.Messages.Message(name, data);
        }

        public static byte[] GetPong()
        {
            Pong pong = new Pong();
            return pong.ToPacketData();
        }
    }

    public static class MessageExtension
    {
        public static byte[] ToPacketData(this IMessage obj)
        {
            if (Gogs.Messages.Instance().Encode(obj, out Gogs.Packet packet))
            {
                return packet.ToByteArray();
            }

            return new byte[] { };
        }
    }


    public static class GEvents
    {
        public static void OnPing(Action<string, Ping> action)
        {
            Gogs.EventsManager.AddListener<Ping>(action);
        }
        
        public static void OnBindUser(Action<string, BindUser> action)
        {
            Gogs.EventsManager.AddListener<BindUser>(action);
        }
        
        public static void OnBindSuccess(Action<string, BindSuccess> action)
        {
            Gogs.EventsManager.AddListener<BindSuccess>(action);
        }
        
    }
}
