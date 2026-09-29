import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createService, type CreateServicePayload } from '@/services/service';

export function useCreateService() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateServicePayload) => createService(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services'] });
    },
  });
}

export type { CreateServicePayload };
