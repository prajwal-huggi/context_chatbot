import React, { useState } from "react";
import axios from "axios";

function ChatPage() {
  const [question, setQuestion] = useState("");
  const [chatHistory, setChatHistory] = useState([]);
  const [loading, setLoading] = useState(false);

const backendUrl = `http://${window.location.hostname}:8080`;
 //import.meta.env.VITE_GO_BACKEND_URL; // Ensure this is in your .env

  // Handle sending a question
  const handleSend = async () => {
    if (!question.trim()) return;

    setLoading(true);
    try {
      const response = await axios.post(`${backendUrl}/api/answer`, {
        question: question.trim(),
      });

      const { mode, answer } = response.data;
      setChatHistory([...chatHistory, { question, mode, answer }]);
      setQuestion(""); // clear input
    } catch (error) {
      console.error(error);
      setChatHistory([...chatHistory, { question, mode: "error", answer: "Failed to get response." }]);
    } finally {
      setLoading(false);
    }
  };

  // Handle reset
  const handleReset = async () => {
    setChatHistory([]); // clear chat

    try {
      const response = await axios.post(`${backendUrl}/api/reset`);
      alert(response.data.message || "Chat reset successfully!");
    } catch (error) {
      console.error(error);
      alert("Failed to reset chat.");
    }
  };

  return (
    <div style={{ maxWidth: "600px", margin: "50px auto", display: "flex", flexDirection: "column", gap: "20px" }}>
      <h2>Chat Page</h2>

      <div style={{ display: "flex", flexDirection: "column", gap: "10px", border: "1px solid #ccc", padding: "10px", borderRadius: "8px", minHeight: "300px" }}>
        {chatHistory.length === 0 ? (
          <p style={{ color: "#888" }}>No messages yet. Start the conversation!</p>
        ) : (
          chatHistory.map((chat, index) => (
            <div key={index} style={{ marginBottom: "10px" }}>
              <p><strong>You:</strong> {chat.question}</p>
              <p><strong>AI ({chat.mode}):</strong> {chat.answer}</p>
            </div>
          ))
        )}
      </div>

      <input
        type="text"
        placeholder="Type your question..."
        value={question}
        onChange={(e) => setQuestion(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleSend()}
        style={{ padding: "10px", fontSize: "16px", borderRadius: "6px", border: "1px solid #ccc" }}
      />

      <div style={{ display: "flex", gap: "10px" }}>
        <button onClick={handleSend} disabled={loading} style={{ padding: "10px 20px", fontSize: "16px" }}>
          {loading ? "Sending..." : "Send"}
        </button>
        <button onClick={handleReset} style={{ padding: "10px 20px", fontSize: "16px", backgroundColor: "#f44336", color: "#fff" }}>
          Reset
        </button>
      </div>
    </div>
  );
}

export default ChatPage;
