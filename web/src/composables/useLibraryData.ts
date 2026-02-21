import { ref, onMounted } from "vue";
import { fetchLibrary, fetchBasePath } from "@/lib/api";
import type { LibraryResponse } from "@/types";

export function useLibraryData() {
    const library = ref<LibraryResponse | null>(null);
    const loading = ref(true);
    const error = ref<string | null>(null);
    const basePath = ref("");

    async function loadLibrary() {
        loading.value = true;
        error.value = null;
        try {
            const [lib, bp] = await Promise.all([fetchLibrary(), fetchBasePath()]);
            library.value = lib;
            basePath.value = bp;
        } catch (e: any) {
            error.value = e.message || "Failed to load library";
        } finally {
            loading.value = false;
        }
    }

    async function refreshLibrary() {
        await loadLibrary();
    }

    onMounted(() => {
        loadLibrary();
    });

    return {
        library,
        loading,
        error,
        basePath,
        refreshLibrary,
    };
}
