using System;
using System.Collections.Generic;
using Google.Protobuf;

namespace Gogs
{
    public static class EventsManager
    {
        static readonly Dictionary<Type, Action<String, IMessage>> s_Events = new Dictionary<Type, Action<String, IMessage>>();

        static readonly Dictionary<Delegate, Action<String, IMessage>> s_EventLookups =
           new Dictionary<Delegate, Action<String, IMessage>>();

        public static void AddListener<T>(Action<String, T> evt) where T : IMessage
        {
            if (!s_EventLookups.ContainsKey(evt))
            {
                Action<String, IMessage> newAction = (s, e) => evt(s, (T)e);
                s_EventLookups[evt] = newAction;

                if (s_Events.TryGetValue(typeof(T), out Action<String, IMessage> internalAction))
                    s_Events[typeof(T)] = internalAction += newAction;
                else
                    s_Events[typeof(T)] = newAction;
            }
        }

        public static void RemoveListener<T>(Action<T> evt) where T : IMessage
        {
            if (s_EventLookups.TryGetValue(evt, out var action))
            {
                if (s_Events.TryGetValue(typeof(T), out var tempAction))
                {
                    tempAction -= action;
                    if (tempAction == null)
                        s_Events.Remove(typeof(T));
                    else
                        s_Events[typeof(T)] = tempAction;
                }

                s_EventLookups.Remove(evt);
            }
        }

        public static void Broadcast(string name, IMessage evt)
        {
            if (s_Events.TryGetValue(evt.GetType(), out var action))
                action.Invoke(name, evt);
        }

        public static void Clear()
        {
            s_Events.Clear();
            s_EventLookups.Clear();
        }
    }
}