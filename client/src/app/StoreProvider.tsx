'use client';
import { useEffect, useState } from 'react';
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { makePersistor, makeStore } from '@/store/store';

export default function StoreProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  // Lazy initializer chạy đúng 1 lần — tạo store + persistor đi kèm.
  const [{ store, persistor }] = useState(() => {
    const store = makeStore();
    return { store, persistor: makePersistor(store) };
  });

  // Expose store cho apiClient có thể dispatch setAccessToken sau khi refresh
  useEffect(() => {
    (window as any).__REDUX_STORE__ = store;
    return () => {
      if ((window as any).__REDUX_STORE__ === store) delete (window as any).__REDUX_STORE__;
    };
  }, [store]);

  return (
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        {children}
      </PersistGate>
    </Provider>
  );
}
