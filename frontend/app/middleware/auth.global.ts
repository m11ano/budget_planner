export default defineNuxtRouteMiddleware((to) => {
    const { $authData } = useNuxtApp();
    const routeName = to.name as string;
    if (routeName === 'index') return;

    const guestRoutes = ['registration'];

    if ($authData.isAuth) {
        if (guestRoutes.includes(routeName)) {
            return navigateTo('/');
        }
    } else {
        if (!guestRoutes.includes(routeName)) {
            return navigateTo('/');
        }
    }
});
