import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { Navbar } from "./components/Navbar";
import { Home } from "./pages/Home";
import { Cinemas } from "./pages/Cinemas";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { MovieDetail } from "./pages/MovieDetail";
import { CinemaDetail } from "./pages/CinemaDetail";
import { HallDetail } from "./pages/HallDetail";
import { Booking } from "./pages/Booking";
import { BookingConfirm } from "./pages/BookingConfirm";
import { Bookings } from "./pages/Bookings";
import { Tickets } from "./pages/Tickets";
import { Profile } from "./pages/Profile";
import { ProtectedRoute } from "./pages/ProtectedRoute";

function App() {
  return (
    <Router>
      <AuthProvider>
        <div className="min-h-screen bg-dark-950">
          <Navbar />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/cinemas" element={<Cinemas />} />
            <Route path="/cinema/:id" element={<CinemaDetail />} />
            <Route
              path="/cinema/:cinemaId/hall/:hallId"
              element={<HallDetail />}
            />
            <Route path="/movie/:id" element={<MovieDetail />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route
              path="/booking/:sessionId"
              element={
                <ProtectedRoute>
                  <Booking />
                </ProtectedRoute>
              }
            />
            <Route
              path="/bookings/:bookingId/purchase"
              element={
                <ProtectedRoute>
                  <BookingConfirm />
                </ProtectedRoute>
              }
            />
            <Route
              path="/bookings"
              element={
                <ProtectedRoute>
                  <Bookings />
                </ProtectedRoute>
              }
            />
            <Route
              path="/tickets"
              element={
                <ProtectedRoute>
                  <Tickets />
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile"
              element={
                <ProtectedRoute>
                  <Profile />
                </ProtectedRoute>
              }
            />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </AuthProvider>
    </Router>
  );
}

export default App;
