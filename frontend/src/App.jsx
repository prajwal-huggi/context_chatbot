import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import ChatPage from './pages/ChatPage'
import Home from './pages/Home'
import UploadPage from './pages/UploadPage'

function App() {

  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/upload" element= {<UploadPage />} />
          <Route path="/chat" element={<ChatPage />} />
        </Routes>

      </Router>
    </>
  )
}

export default App
