import { useState, useEffect, useCallback } from "react";
import { IDB } from "@/utils/indexedDB";

export function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T | ((val: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      const item = window.localStorage.getItem(key);
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.error(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  });

  const setValue = (value: T | ((val: T) => T)) => {
    try {
      const valueToStore =
        value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      window.localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.error(`Error setting localStorage key "${key}":`, error);
    }
  };

  useEffect(() => {
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === key && e.newValue) {
        try {
          setStoredValue(JSON.parse(e.newValue));
        } catch (error) {
          console.error(
            `Error parsing localStorage value for key "${key}":`,
            error
          );
        }
      }
    };

    window.addEventListener("storage", handleStorageChange);
    return () => window.removeEventListener("storage", handleStorageChange);
  }, [key]);

  return [storedValue, setValue];
}

export function useIndexedDB<T>(
  key: string,
  initialValue: T
): [T, (value: T | ((val: T) => T)) => void, boolean] {
  const [storedValue, setStoredValue] = useState<T>(initialValue);
  const [isLoading, setIsLoading] = useState(true);

  // Load initial value from IndexedDB
  useEffect(() => {
    let isMounted = true;

    const loadInitialValue = async () => {
      try {
        setIsLoading(true);
        const item = await IDB.getItem<T>(key);

        if (isMounted) {
          setStoredValue(item !== null ? item : initialValue);
          setIsLoading(false);
        }
      } catch (error) {
        console.error(`Error reading IndexedDB key "${key}":`, error);
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    loadInitialValue();

    return () => {
      isMounted = false;
    };
  }, [key, initialValue]);

  // Set value function
  const setValue = useCallback(
    async (value: T | ((val: T) => T)) => {
      try {
        const valueToStore =
          value instanceof Function ? value(storedValue) : value;
        setStoredValue(valueToStore);
        await IDB.setItem(key, valueToStore);
      } catch (error) {
        console.error(`Error setting IndexedDB key "${key}":`, error);
      }
    },
    [key, storedValue]
  );

  return [storedValue, setValue, isLoading];
}
