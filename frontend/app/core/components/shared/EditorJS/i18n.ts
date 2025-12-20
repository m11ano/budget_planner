import type { I18nConfig } from '@editorjs/editorjs';

export const i18n: I18nConfig = {
    messages: {
        ui: {
            blockTunes: {
                toggler: {
                    'Click to tune': 'Нажмите, чтобы настроить',
                    'or drag to move': 'или перетащите',
                },
            },
            inlineToolbar: {
                converter: {
                    'Convert to': 'Конвертировать в',
                },
            },
            toolbar: {
                toolbox: {
                    Add: 'Добавить',
                },
            },
            popover: {
                Filter: 'Фильтр',
            },
        },
        toolNames: {
            Text: 'Параграф',
            Heading: 'Заголовок',
            List: 'Список',
            Warning: 'Примечание',
            Checklist: 'Чеклист',
            Quote: 'Цитата',
            Code: 'Код',
            Delimiter: 'Разделитель',
            'Raw HTML': 'HTML-фрагмент',
            Table: 'Таблица',
            Link: 'Ссылка',
            Marker: 'Маркер',
            Bold: 'Полужирный',
            Italic: 'Курсив',
            InlineCode: 'Моноширинный',
            Image: 'Изображение',
            Underline: 'Подчеркнутый',
            'Ordered List': 'Нумерированный список',
            'Unordered List': 'Маркерный список',
        },
        tools: {
            image: {
                'Select an Image': 'Загрузить изображение',
                Delete: 'Удалить',
                'Gallery caption': 'Подпись',
            },
            link: {
                Save: 'Сохранить',
                'Select target': 'Выберите цель',
                'Select rel': 'Выберите rel',
                'The URL is not valid.': 'Некорректный адрес ссылки',
            },
            list: {
                Unordered: 'Маркированный',
                Ordered: 'Нумерованный',
            },
            gallery: {
                'Select an Image': 'Загрузить изображение',
                Delete: 'Удалить',
                'Gallery caption': 'Подпись',
            },
        },
        blockTunes: {
            delete: {
                Delete: 'Удалить',
                'Click to delete': 'Подтвердить удаление',
            },
            moveUp: {
                'Move up': 'Переместить вверх',
            },
            moveDown: {
                'Move down': 'Переместить вниз',
            },
        },
    },
};
