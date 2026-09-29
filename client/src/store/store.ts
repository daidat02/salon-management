import { combineReducers, configureStore } from '@reduxjs/toolkit';
import {
  FLUSH,
  PAUSE,
  PERSIST,
  PURGE,
  REGISTER,
  REHYDRATE,
  persistReducer,
  persistStore,
} from 'redux-persist';
import createWebStorage from 'redux-persist/lib/storage/createWebStorage';
import authReducer from './slices/authSlice';

// redux-persist storage chạm window.localStorage — dùng noop khi render
// phía server để makeStore không crash ngoài client component.
const createNoopStorage = () => ({
  getItem: () => Promise.resolve(null),
  setItem: () => Promise.resolve(),
  removeItem: () => Promise.resolve(),
});

const storage =
  typeof window !== 'undefined'
    ? createWebStorage('local')
    : createNoopStorage();

// Chỉ persist user (profile). accessToken/status/error giữ trong memory:
// token trong localStorage là mồi cho XSS; backend chưa có /refresh nên
// reload trang sẽ cần đăng nhập lại để lấy token mới.
const authPersistConfig = {
  key: 'auth',
  storage,
  blacklist: ['accessToken', 'status', 'error'],
};

const persistedAuthReducer = persistReducer(authPersistConfig, authReducer);

const rootReducer = combineReducers({
  auth: authReducer,
});

const persistedReducer = combineReducers({
  auth: persistedAuthReducer,
});

export const makeStore = () => {
  return configureStore({
    reducer: persistedReducer,
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware({
        serializableCheck: {
          ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
        },
      }),
  });
};

// Infer the type of makeStore
export type AppStore = ReturnType<typeof makeStore>;
// Infer the `RootState` and `AppDispatch` types from the root reducer
export type RootState = ReturnType<typeof rootReducer>;
export type AppDispatch = AppStore['dispatch'];

// Tạo persistor đi kèm 1 store instance (gọi 1 lần trong StoreProvider).
export const makePersistor = (store: AppStore) => persistStore(store);
