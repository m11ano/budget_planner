export function showErrors(errors: string[] = ['Неизвестная ошибка']) {
    useToast().add({
        title: 'Возникли ошибки',
        description: errors.join('\n'),
        color: 'error',
        icon: 'i-lucide-alert-triangle',
    });
}

export function showSuccess(msg: string = 'Данные сохранены') {
    useToast().add({
        title: 'Успех',
        description: msg,
        color: 'success',
        icon: 'i-lucide-check-circle',
    });
}
