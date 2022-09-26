using System;
using Gogs;
using Google.Protobuf;
namespace {{.Package | CamelCase }}
{
    public static class GInit
    {
        public const int SERVICE_TYPE = 2;

        public static void Init()
        {
            {{range .Components}}{{range .Fields}}
            Dispatch.AddField("{{.Name}}", typeof({{.Name}}), Packet.CreateAction(SERVICE_TYPE, {{.ComponentIndex}}, {{.Index}}));  // {{.Action16}} {{.Action10}} 
            {{end}}{{end}}
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
        {{range .Components}}{{range .Fields}}
        public static void On{{.Name}}(Action<string, {{.Name}}> action)
        {
            Gogs.EventsManager.AddListener<{{.Name}}>(action);
        }
        {{end}}{{end}}
    }
}
