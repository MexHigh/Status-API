import Card from "./Card";
import WidthContainer from "./WidthContainer";

export default function LoginForm() {
    return (
        <WidthContainer className="my-8">
            <Card>
                <div className="flex flex-col space-y-4">
                    <label for="password" className="text-xl">Message API key</label>
                    <input 
                        type="password"
                        className="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-0 focus:border-green-500"
                        placeholder="Enter your Message API key"
                    />
                    <button 
                        type="submit"
                        className="w-full px-4 py-2 bg-green-500 text-white font-semibold rounded-md shadow-sm hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-1"
                    >
                        Login
                    </button>
                </div>
            </Card>
        </WidthContainer>
    )
}