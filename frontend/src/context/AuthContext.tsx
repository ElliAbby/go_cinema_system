import React, { createContext, useContext, useState, useEffect } from "react";
import { User } from "../types";

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const loadStoredUser = (): User | null => {
  const token = localStorage.getItem("auth_token");
  const savedUser = localStorage.getItem("user");

  if (!token || !savedUser) {
    return null;
  }

  try {
    return JSON.parse(savedUser) as User;
  } catch (error) {
    console.error("Failed to parse user data", error);
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user");
    return null;
  }
};

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(loadStoredUser);

  useEffect(() => {
    if (!user) {
      return;
    }

    localStorage.setItem("user", JSON.stringify(user));
  }, [user]);

  const login = (token: string, userData: User) => {
    localStorage.setItem("auth_token", token);
    localStorage.setItem("user", JSON.stringify(userData));
    setUser(userData);
  };

  const logout = () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user");
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!localStorage.getItem("auth_token") && !!user,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
};
