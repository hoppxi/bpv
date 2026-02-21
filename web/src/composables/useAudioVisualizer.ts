import { ref, watch, onUnmounted } from "vue";

export function useAudioVisualizer(getAnalyser: () => AnalyserNode | null) {
    const frequencyData = ref<Uint8Array>(new Uint8Array(256));
    const timeDomainData = ref<Uint8Array>(new Uint8Array(256));
    const isActive = ref(false);

    let animationId: number | null = null;

    function start() {
        if (isActive.value) return;
        isActive.value = true;
        update();
    }

    function stop() {
        isActive.value = false;
        if (animationId !== null) {
            cancelAnimationFrame(animationId);
            animationId = null;
        }
    }

    function update() {
        if (!isActive.value) return;

        const analyser = getAnalyser();
        if (!analyser) {
            animationId = requestAnimationFrame(update);
            return;
        }

        const freq = new Uint8Array(analyser.frequencyBinCount);
        const time = new Uint8Array(analyser.frequencyBinCount);

        analyser.getByteFrequencyData(freq);
        analyser.getByteTimeDomainData(time);

        frequencyData.value = freq;
        timeDomainData.value = time;

        animationId = requestAnimationFrame(update);
    }

    onUnmounted(() => {
        stop();
    });

    return {
        frequencyData,
        timeDomainData,
        isActive,
        start,
        stop,
    };
}
