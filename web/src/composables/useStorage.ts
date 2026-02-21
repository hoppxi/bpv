import { ref, watch, onMounted, type Ref } from "vue";
import { fetchSettings, saveSettings, type SettingsState } from "@/lib/api";

type SupportedKey =
  | "bpv-visualizer-type"
  | "bpv-show-visualizer"
  | "bpv-auto-play"
  | "bpv-crossfade"
  | "bpv-gapless"
  | "bpv-eq";

let cachedSettings: SettingsState | null = null;
let loadPromise: Promise<SettingsState> | null = null;
let saveTimer: number | null = null;

async function ensureLoaded(): Promise<SettingsState> {
  if (cachedSettings) return cachedSettings;
  if (!loadPromise) {
    loadPromise = fetchSettings()
      .catch(() => ({} as SettingsState))
      .then((s) => {
        cachedSettings = s || {};
        return cachedSettings;
      });
  }
  return loadPromise;
}

function scheduleSave() {
  if (!cachedSettings) return;
  if (saveTimer) window.clearTimeout(saveTimer);
  saveTimer = window.setTimeout(() => {
    if (!cachedSettings) return;
    saveSettings(cachedSettings).catch(() => {});
  }, 250);
}

function getFromSettings<T>(key: SupportedKey, settings: SettingsState): T | undefined {
  switch (key) {
    case "bpv-visualizer-type":
      return settings.visualizer_type as T | undefined;
    case "bpv-show-visualizer":
      return settings.show_visualizer as T | undefined;
    case "bpv-auto-play":
      return settings.auto_play as T | undefined;
    case "bpv-crossfade":
      return settings.crossfade as T | undefined;
    case "bpv-gapless":
      return settings.gapless as T | undefined;
    case "bpv-eq":
      return (
        settings.eq_bass !== undefined ||
        settings.eq_mid !== undefined ||
        settings.eq_treble !== undefined ||
        settings.eq_enabled !== undefined
          ? ({
              bass: settings.eq_bass ?? 0,
              mid: settings.eq_mid ?? 0,
              treble: settings.eq_treble ?? 0,
              enabled: settings.eq_enabled ?? false,
            } as T)
          : undefined
      );
  }
}

function setInSettings<T>(key: SupportedKey, settings: SettingsState, value: T) {
  switch (key) {
    case "bpv-visualizer-type":
      settings.visualizer_type = value as any;
      return;
    case "bpv-show-visualizer":
      settings.show_visualizer = value as any;
      return;
    case "bpv-auto-play":
      settings.auto_play = value as any;
      return;
    case "bpv-crossfade":
      settings.crossfade = value as any;
      return;
    case "bpv-gapless":
      settings.gapless = value as any;
      return;
    case "bpv-eq": {
      const v = value as any;
      settings.eq_bass = typeof v?.bass === "number" ? v.bass : 0;
      settings.eq_mid = typeof v?.mid === "number" ? v.mid : 0;
      settings.eq_treble = typeof v?.treble === "number" ? v.treble : 0;
      settings.eq_enabled = !!v?.enabled;
      return;
    }
  }
}

export function useStorage<T>(key: string, defaultValue: T): Ref<T> {
  const data = ref<T>(structuredClone(defaultValue)) as Ref<T>;
  let initialized = false;

  onMounted(async () => {
    const k = key as SupportedKey;
    const settings = await ensureLoaded();
    const stored = getFromSettings<T>(k, settings);
    if (stored !== undefined) {
      data.value = stored;
    }
    initialized = true;
  });

  watch(
    data,
    (newValue) => {
      if (!initialized) return;
      const k = key as SupportedKey;
      if (!cachedSettings) cachedSettings = {};
      setInSettings<T>(k, cachedSettings, newValue);
      scheduleSave();
    },
    { deep: true }
  );

  return data;
}
