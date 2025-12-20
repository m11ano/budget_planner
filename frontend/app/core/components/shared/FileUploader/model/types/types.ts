import type { IFile } from '~/shared/types/files';

export interface UploadedFile {
    targets: Record<string, IFile>;
}
