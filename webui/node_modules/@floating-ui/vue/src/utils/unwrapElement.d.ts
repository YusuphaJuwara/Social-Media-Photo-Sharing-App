import type { ComponentPublicInstance } from 'vue-demi';
import type { MaybeElement } from '../types';
export declare function unwrapElement<T>(element: MaybeElement<T>): Exclude<T, ComponentPublicInstance> | null | undefined;
