export default function TextInput({ label, ...props }) {
    return (
        <div>
            <label htmlFor={label} className="text-lg">{label}</label>
            <input
                id={label}
                className="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-0 focus:border-green-500"
                {...props}
            />
        </div>
    )
}