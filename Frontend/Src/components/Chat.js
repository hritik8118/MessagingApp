import React, { useEffect, useRef, useState } from "react";

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [receiver, setReceiver] = useState("");
  const [content, setContent] = useState("");
  const ws = useRef(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    ws.current = new WebSocket(`ws://localhost:8080/chat/ws?token=${token}`);
    ws.current.onmessage = (event) => {
      setMessages((msgs) => [...msgs, JSON.parse(event.data)]);
    };
    ws.current.onclose = () => alert("WebSocket closed");
    return () => ws.current.close();
  }, []);

  const sendMessage = (e) => {
    e.preventDefault();
    ws.current.send(JSON.stringify({ receiver, content }));
    setContent("");
  };

  return (
    <div>
      <h2>Chat</h2>
      <form onSubmit={sendMessage}>
        <input placeholder="Receiver" value={receiver} onChange={e => setReceiver(e.target.value)} />
        <input placeholder="Message" value={content} onChange={e => setContent(e.target.value)} />
        <button type="submit">Send</button>
      </form>
      <ul>
        {messages.map((msg, i) => (
          <li key={i}><b>{msg.sender}:</b> {msg.content}</li>
        ))}
      </ul>
    </div>
  );
}
