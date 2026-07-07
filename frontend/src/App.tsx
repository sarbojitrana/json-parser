import Home from "./pages/Home"
import { ToastContainer } from "react-toastify"
import "react-toastify/dist/ReactToastify.css"

function App() {
    return (
        <>
            <Home />
            <ToastContainer 
                position="top-right" 
                autoClose={4000} 
                hideProgressBar={true} 
                newestOnTop={false}
                closeOnClick
                pauseOnHover
                theme="colored"
            />
        </>
    )
}

export default App