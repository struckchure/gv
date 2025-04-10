import * as React from 'https://esm.sh/@types/react@~19.1.0/index.d.ts';
import { R as RouterProviderProps$1 } from './fog-of-war-1hWhK5ey.d.mts';
import { R as RouterInit } from './route-data-5OzAzQtT.d.mts';

type RouterProviderProps = Omit<RouterProviderProps$1, "flushSync">;
declare function RouterProvider(props: Omit<RouterProviderProps, "flushSync">): React.JSX.Element;

interface HydratedRouterProps {
    /**
     * Context object to passed through to `createBrowserRouter` and made available
     * to `clientLoader`/`clientActon` functions
     */
    unstable_getContext?: RouterInit["unstable_getContext"];
}
/**
 * Framework-mode router component to be used in `entry.client.tsx` to hydrate a
 * router from a `ServerRouter`
 *
 * @category Component Routers
 */
declare function HydratedRouter(props: HydratedRouterProps): React.JSX.Element;

export { HydratedRouter, RouterProvider, type RouterProviderProps };