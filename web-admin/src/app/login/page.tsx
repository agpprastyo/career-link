import {LoginForm} from "@/components/login-form"
import {ModeToggle} from "@/components/theme-toggle";

export default function Page() {
    return (
        <div className="relative flex min-h-screen w-full items-center justify-center p-6 md:p-10 overflow-hidden">
            <div className="absolute top-4 right-4">
                <ModeToggle/>
            </div>
            <div className="w-full max-w-sm">
                <LoginForm/>
            </div>
        </div>
    )
}