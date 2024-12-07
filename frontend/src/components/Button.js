export const primaryCss = "bg-green-500 text-white hover:bg-green-600"
export const secondaryCss = "ring-1 ring-gray-200 hover:bg-gray-100"

export default function Button({ text, wFull = true, primary = false, ...props }) {
    return (
        <button
            className={`px-4 py-2 rounded-md shadow-sm ${primary ? primaryCss : secondaryCss} ${wFull && 'w-full'}`}
            {...props}
        >
            {text}
        </button>
    )
}
