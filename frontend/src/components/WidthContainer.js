export default function WidthContainer({ children, className }) {
    let base = "mx-auto max-w-5xl "
    return (
        <div className={base + className}>
            { children }
        </div>
    )
}