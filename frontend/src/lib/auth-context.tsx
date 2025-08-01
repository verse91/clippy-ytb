"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { User } from "@supabase/supabase-js";
import { supabase } from "./supabaseClients";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  signOut: () => Promise<void>;
  getUserCredits: (userId: string) => Promise<number>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Get initial session
    const getSession = async () => {
      const {
        data: { session },
      } = await supabase.auth.getSession();
      setUser(session?.user ?? null);
      setLoading(false);
    };

    getSession();

    // Listen for auth changes
    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange(async (event, session) => {
      setUser(session?.user ?? null);
      setLoading(false);
    });

    return () => subscription.unsubscribe();
  }, []);

  const signOut = async () => {
    await supabase.auth.signOut();
  };

  const getUserCredits = async (userId: string): Promise<number> => {
    try {
      const response = await fetch(`/api/v1/users/${userId}/credits`, {
        headers: {
          "X-User-ID": userId,
        },
      });

      if (!response.ok) {
        console.error("Failed to fetch user credits");
        return 0;
      }

      const data = await response.json();
      return data.data?.credits || 0;
    } catch (error) {
      console.error("Error fetching user credits:", error);
      return 0;
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, signOut, getUserCredits }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
