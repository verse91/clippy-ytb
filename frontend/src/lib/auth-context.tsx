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
      try {
        const {
          data: { session },
        } = await supabase.auth.getSession();
        setUser(session?.user ?? null);
      } catch (error) {
        console.error("Failed to get initial session:", error);
        setUser(null);
      } finally {
        setLoading(false);
      }
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
    try {
      await supabase.auth.signOut();
    } catch (error) {
      console.error("Failed to sign out:", error);
      throw error; // Re-throw to allow caller to handle
    }
  };

  const getUserCredits = async (userId: string): Promise<number> => {
    if (!userId?.trim()) {
      console.error("Invalid userId provided");
      return 0;
    }
    try {
      const response = await fetch(`/api/v1/users/${userId}/credits`);

      if (!response.ok) {
        console.error(`Failed to fetch user credits: ${response.status} ${response.statusText}`);
        return 0;
      }

      const data = await response.json();
      return typeof data.credits === 'number' ? data.credits : (data.data?.credits || 0);
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
