import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import type { UseQueryResult } from "@tanstack/react-query"
import axiosInstance from "../utils/axiosConfig"
import { parseApiError } from "../utils/parseApiError"
import type { Student, CreateStudentPayload, UpdateStudentPayload } from "../types/students"

const API_URL = import.meta.env.VITE_API_URL as string + "/api"

export const useGetStudents = (): UseQueryResult<Student[], Error> =>
    useQuery<Student[], Error>({
        queryKey: ["students"],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: Student[] }>(`${API_URL}/students`)
                return data.data ?? []
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        refetchInterval: 5000,
    })

export const useCreateStudent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (payload: CreateStudentPayload) => {
            try {
                const { data } = await axiosInstance.post<{ data: Student }>(`${API_URL}/student`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["students"] }),
    })
}

export const useUpdateStudent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async ({ id, payload }: { id: number; payload: UpdateStudentPayload }) => {
            try {
                const { data } = await axiosInstance.put<{ data: Student }>(`${API_URL}/student/${id}`, payload)
                return data.data
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["students"] }),
    })
}

export const useDeleteStudent = () => {
    const queryClient = useQueryClient()
    return useMutation({
        mutationFn: async (id: number) => {
            try {
                await axiosInstance.delete(`${API_URL}/student/${id}`)
            } catch (err) {
                throw new Error(parseApiError(err))
            }
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["students"] }),
    })
}
