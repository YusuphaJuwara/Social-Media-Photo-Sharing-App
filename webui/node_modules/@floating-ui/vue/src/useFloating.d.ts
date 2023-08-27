import type { FloatingElement, ReferenceElement } from '@floating-ui/dom';
import type { Ref } from 'vue-demi';
import type { MaybeElement, UseFloatingOptions, UseFloatingReturn } from './types';
/**
 * Computes the `x` and `y` coordinates that will place the floating element next to a reference element when it is given a certain CSS positioning strategy.
 * @param reference The reference template ref.
 * @param floating The floating template ref.
 * @param options The floating options.
 * @see https://floating-ui.com/docs/vue
 */
export declare function useFloating<T extends ReferenceElement = ReferenceElement>(reference: Readonly<Ref<MaybeElement<T>>>, floating: Readonly<Ref<MaybeElement<FloatingElement>>>, options?: UseFloatingOptions<T>): UseFloatingReturn;
