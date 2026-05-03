import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { UseQueryResult } from '@tanstack/react-query'
import axiosInstance from '../utils/axiosConfig'
import { parseApiError } from '../utils/parseApiError'
import type { Parent, CreateParentPayload, UpdateParentPayload } from '../types/parents'

const API_URL = import.meta.env.VITE_API_URL as string + '/api'

export const useGetParents = (): UseQueryResult<Parent[], Error> =>
    useQuery<Parent[], Error>({
        queryKey: ['parents'],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Parent[] }>(`${API_URL}/parents`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateParent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateParentPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Parent }>(`${API_URL}/parent`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['parents'] }),
    })
}

export const useUpdateParent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async ({ id, payload }: { id: number; payload: UpdateParentPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: Parent }>(`${API_URL}/parent/${id}`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['parents'] }),
    })
}

export const useDeleteParent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/parent/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ['parents'] }),
    })
}

export const useGetStudentsForParent = (parentId: number) =>
    useQuery({
        queryKey: ['parent-students', parentId],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: { student_id: number; first_name: string; last_name: string }[] }>(
                    `${API_URL}/student-parent/parent/${parentId}/students`
                )
                return data.data ?? []
            } catch {
                return []
            }
        },
        enabled: parentId > 0,
    })

export const linkParentToStudent = async (studentId: number, parentId: number) =>
    axiosInstance.post(`${API_URL}/student-parent/student/${studentId}/parent/${parentId}`)

export const unlinkParentFromStudent = async (studentId: number, parentId: number) =>
    axiosInstance.delete(`${API_URL}/student-parent/student/${studentId}/parent/${parentId}`)

