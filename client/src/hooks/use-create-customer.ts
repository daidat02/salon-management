import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createCustomer, type CreateCustomerPayload } from '@/services/customer';
import { getApiErrorMessage } from '@/services/apiClient';

export function useCreateCustomer() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateCustomerPayload) => createCustomer(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] });
    },
  });
}

export type { CreateCustomerPayload };
