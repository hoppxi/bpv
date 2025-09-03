export class IDB {
  private static db: IDBDatabase | null = null;
  private static dbName = "bpvCache";
  private static storeName = "keyValueStore";
  private static isInitialized = false;

  private static async initialize(): Promise<void> {
    if (this.isInitialized) return;

    return new Promise((resolve, reject) => {
      const request = indexedDB.open(this.dbName, 1);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        this.db = request.result;
        this.isInitialized = true;
        resolve();
      };

      request.onupgradeneeded = (event: IDBVersionChangeEvent) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(this.storeName)) {
          db.createObjectStore(this.storeName);
        }
      };
    });
  }

  static async setItem<T>(key: string, value: T): Promise<void> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readwrite");
      const store = transaction.objectStore(this.storeName);
      const request = store.put(value, key);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve();
    });
  }

  static async getItem<T>(key: string): Promise<T | null> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.get(key);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result || null);
    });
  }

  static async removeItem(key: string): Promise<void> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readwrite");
      const store = transaction.objectStore(this.storeName);
      const request = store.delete(key);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve();
    });
  }

  static async clear(): Promise<void> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readwrite");
      const store = transaction.objectStore(this.storeName);
      const request = store.clear();

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve();
    });
  }

  static async key(index: number): Promise<string | null> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.getAllKeys();

      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        const keys = request.result as string[];
        resolve(keys[index] || null);
      };
    });
  }

  static async keys(): Promise<string[]> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.getAllKeys();

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result as string[]);
    });
  }

  static async hasItem(key: string): Promise<boolean> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.getKey(key);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result !== undefined);
    });
  }

  static async length(): Promise<number> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.count();

      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result);
    });
  }

  static async getAll<T>(): Promise<{ [key: string]: T }> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.getAll();

      request.onerror = () => reject(request.error);
      request.onsuccess = async () => {
        const values = request.result as T[];
        const keys = await this.keys();
        const result: { [key: string]: T } = {};

        keys.forEach((key, index) => {
          result[key] = values[index];
        });

        resolve(result);
      };
    });
  }

  static async getAllEntries<T>(): Promise<[string, T][]> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.getAll();

      request.onerror = () => reject(request.error);
      request.onsuccess = async () => {
        const values = request.result as T[];
        const keys = await this.keys();
        const entries: [string, T][] = keys.map((key, index) => [
          key,
          values[index],
        ]);
        resolve(entries);
      };
    });
  }

  static async getMultiple<T>(
    keys: string[]
  ): Promise<{ [key: string]: T | null }> {
    await this.initialize();

    const results: { [key: string]: T | null } = {};

    // Get all items in parallel for better performance
    await Promise.all(
      keys.map(async (key) => {
        try {
          results[key] = await this.getItem<T>(key);
        } catch (error) {
          console.error(`Error getting key "${key}":`, error);
          results[key] = null;
        }
      })
    );

    return results;
  }

  static async getStorageInfo(): Promise<{
    used: number;
    quota: number;
    usedFormatted: string;
    quotaFormatted: string;
    percentage: number;
  }> {
    try {
      if (!navigator.storage || !navigator.storage.estimate) {
        throw new Error("Storage Estimation API not supported");
      }

      const estimate = await navigator.storage.estimate();

      return {
        used: estimate.usage || 0,
        quota: estimate.quota || 0,
        usedFormatted: this.formatBytes(estimate.usage || 0),
        quotaFormatted: this.formatBytes(estimate.quota || 0),
        percentage: estimate.quota
          ? ((estimate.usage || 0) / estimate.quota) * 100
          : 0,
      };
    } catch (error) {
      console.error("Error getting storage info:", error);
      return {
        used: 0,
        quota: 0,
        usedFormatted: "0 Bytes",
        quotaFormatted: "0 Bytes",
        percentage: 0,
      };
    }
  }

  static async estimateDatabaseSize(): Promise<number> {
    await this.initialize();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([this.storeName], "readonly");
      const store = transaction.objectStore(this.storeName);
      const request = store.count();

      request.onerror = () => reject(request.error);
      request.onsuccess = async () => {
        const count = request.result;

        if (count === 0) {
          resolve(0);
          return;
        }

        // Sample a few records to estimate average size
        const sampleSize = Math.min(10, count);
        let totalSampleSize = 0;
        let samplesProcessed = 0;

        const cursorRequest = store.openCursor();

        cursorRequest.onsuccess = function (event) {
          const cursor = (event.target as IDBRequest<IDBCursorWithValue>)
            .result;
          if (cursor && samplesProcessed < sampleSize) {
            const recordSize = JSON.stringify(cursor.value).length;
            totalSampleSize += recordSize;
            samplesProcessed++;
            cursor.continue();
          } else {
            const averageSize =
              samplesProcessed > 0 ? totalSampleSize / samplesProcessed : 0;
            const estimatedSize = averageSize * count;
            resolve(estimatedSize);
          }
        };

        cursorRequest.onerror = () => reject(cursorRequest.error);
      };
    });
  }

  private static formatBytes(bytes: number, decimals: number = 2): string {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  static async getDetailedStorageInfo(): Promise<{
    origin: {
      used: number;
      quota: number;
      usedFormatted: string;
      quotaFormatted: string;
      percentage: number;
    };
    database: {
      estimatedSize: number;
      estimatedSizeFormatted: string;
      recordCount: number;
    };
  }> {
    const [originInfo, dbSize, recordCount] = await Promise.all([
      this.getStorageInfo(),
      this.estimateDatabaseSize(),
      this.length(),
    ]);

    return {
      origin: originInfo,
      database: {
        estimatedSize: dbSize,
        estimatedSizeFormatted: this.formatBytes(dbSize),
        recordCount: recordCount,
      },
    };
  }
}
