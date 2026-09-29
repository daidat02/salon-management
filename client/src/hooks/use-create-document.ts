import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createInventoryDocument, type CreateDocumentPayload } from '@/services/inventory';

export function useCreateDocument() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateDocumentPayload) => createInventoryDocument(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory-documents'] });
      queryClient.invalidateQueries({ queryKey: ['inventory-transactions'] });
      queryClient.invalidateQueries({ queryKey: ['products'] });
    },
  });
}

export type { CreateDocumentPayload };
