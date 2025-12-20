import { useMediaQuery } from '@vueuse/core';

const size = {
    sm: 576,
    md: 768,
    lg: 992,
    xl: 1200,
    xxl: 1600,
};

export function useScreenSizes() {
    const smMore = useMediaQuery(`(min-width: ${size.sm}px)`);
    const smLess = useMediaQuery(`(max-width: ${size.sm - 1}px)`);

    const mdMore = useMediaQuery(`(min-width: ${size.md}px)`);
    const mdLess = useMediaQuery(`(max-width: ${size.md - 1}px)`);

    const lgMore = useMediaQuery(`(min-width: ${size.lg}px)`);
    const lgLess = useMediaQuery(`(max-width: ${size.lg - 1}px)`);

    const xlMore = useMediaQuery(`(min-width: ${size.xl}px)`);
    const xlLess = useMediaQuery(`(max-width: ${size.xl - 1}px)`);

    const xxlMore = useMediaQuery(`(min-width: ${size.xxl}px)`);
    const xxlLess = useMediaQuery(`(max-width: ${size.xxl - 1}px)`);

    return {
        smMore,
        smLess,
        mdMore,
        mdLess,
        lgMore,
        lgLess,
        xlMore,
        xlLess,
        xxlMore,
        xxlLess,
    };
}
