import { ref, computed } from "vue";
import {
    fetchFavorites as apiFetchFavorites,
    addFavorite as apiAddFavorite,
    removeFavorite as apiRemoveFavorite,
} from "@/lib/api";

const favoritesData = ref<Set<string>>(new Set());
const initialized = ref(false);
const isLoading = ref(false);

export function useFavorites() {
    async function loadFavorites() {
        if (initialized.value && !isLoading.value) return;
        isLoading.value = true;
        try {
            const data = await apiFetchFavorites();
            if (Array.isArray(data)) {
                favoritesData.value = new Set(data);
                initialized.value = true;
            }
        } catch (e) {
            console.error("Failed to load favorites:", e);
        } finally {
            isLoading.value = false;
        }
    }

    // Load initially if not loaded
    if (!initialized.value && !isLoading.value) {
        loadFavorites();
    }

    async function toggleFavorite(filePath?: string) {
        if (!filePath) return;

        const isFav = favoritesData.value.has(filePath);

        // Optimistic update
        const newSet = new Set(favoritesData.value);
        if (isFav) {
            newSet.delete(filePath);
        } else {
            newSet.add(filePath);
        }
        favoritesData.value = newSet;

        try {
            if (isFav) {
                await apiRemoveFavorite(filePath);
            } else {
                await apiAddFavorite(filePath);
            }
        } catch (e) {
            console.error("Failed to toggle favorite:", e);
            // Revert on error
            if (isFav) {
                favoritesData.value.add(filePath);
            } else {
                favoritesData.value.delete(filePath);
            }
            favoritesData.value = new Set(favoritesData.value);
        }
    }

    const isFavorite = (filePath?: string | null) => {
        if (!filePath) return false;
        return favoritesData.value.has(filePath);
    };

    return {
        favorites: favoritesData,
        favoritesSet: favoritesData, // Alias for compatibility
        isLoading,
        loadFavorites,
        toggleFavorite,
        isFavorite,
        refreshFavorites: async () => {
            initialized.value = false;
            await loadFavorites();
        },
    };
}
