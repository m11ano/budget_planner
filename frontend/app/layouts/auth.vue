<template>
    <div id="auth_layout">
        <WidgetHeader />
        <div class="content">
            <aside class="side_menu"><WidgetMenu /></aside>
            <main class="main">
                <div class="breadcrumbs"><WidgetBreadcrumb /></div>
                <div class="wrapper"><slot /></div>
            </main>
        </div>
    </div>
</template>
