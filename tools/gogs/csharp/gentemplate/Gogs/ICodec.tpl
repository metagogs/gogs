using Google.Protobuf;
namespace Gogs
{
    public interface IEncode
    {
        public byte[] Encode(IMessage msg);
    }

    public interface IDecode
    {
        public T Decode<T>(byte[] data) where T : IMessage<T>, new();
    }
}