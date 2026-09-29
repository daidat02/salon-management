import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createStaff, type CreateStaffPayload } from '@/services/staff';

export function useCreateStaff() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateStaffPayload) => createStaff(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['staff'] });
    },
  });
}

export type { CreateStaffPayload };
