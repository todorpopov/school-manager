export const inputCls = (hasError: boolean) =>
    `w-full px-3 py-2 text-sm border rounded-md bg-white dark:bg-slate-900 text-slate-800 dark:text-slate-100 focus:outline-none focus:ring-2 transition-colors ${
        hasError
            ? 'border-red-400 focus:ring-red-400'
            : 'border-slate-300 dark:border-slate-600 focus:ring-indigo-400'
    }`

export const primaryBtnCls =
    'px-4 py-2 bg-indigo-500 hover:bg-indigo-600 disabled:opacity-50 text-white text-sm font-medium rounded-md border-none cursor-pointer transition-colors'

