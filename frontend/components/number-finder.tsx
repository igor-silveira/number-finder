"use client"

import {useCallback, useState} from "react"
import {Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle} from "@/components/ui/card"
import {Input} from "@/components/ui/input"
import {Button} from "@/components/ui/button"
import {Alert, AlertDescription, AlertTitle} from "@/components/ui/alert"
import {Loader2} from "lucide-react"
import {Label} from "@/components/ui/label"

interface ApiResponse {
    index: number
    value: number
    is_approximate: boolean
}

export default function NumberFinder() {
    const [number, setNumber] = useState("")
    const [thresholdPercentage, setThresholdPercentage] = useState("")
    const [result, setResult] = useState<ApiResponse | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [isLoading, setIsLoading] = useState(false)
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || ""

    const handleSubmit = useCallback(
        async (e: React.FormEvent) => {
            e.preventDefault()
            setIsLoading(true)
            setError(null)
            setResult(null)

            try {
                let url = `${apiUrl}/api/number/${number}`
                if (thresholdPercentage) {
                    url += `?thresholdPercentage=${thresholdPercentage}`
                }

                const response = await fetch(url)
                const data = await response.json()

                if (!response.ok) {
                    if (response.status === 400) {
                        throw new Error(data.error || "Invalid value parameter")
                    } else if (response.status === 404) {
                        throw new Error(data.error || "No value found within acceptable threshold")
                    } else {
                        throw new Error("An unexpected error occurred")
                    }
                }

                setResult(data as ApiResponse)
            } catch (err) {
                setError(err instanceof Error ? err.message : "An unexpected error occurred")
            } finally {
                setIsLoading(false)
            }
        },
        [number, thresholdPercentage],
    )

    return (
        <Card className="w-[350px]">
            <CardHeader>
                <CardTitle>Number Finder</CardTitle>
                <CardDescription>Find the index of a number in the sequence</CardDescription>
            </CardHeader>
            <CardContent>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="number">Number</Label>
                        <Input
                            id="number"
                            type="number"
                            placeholder="Enter a number"
                            value={number}
                            onChange={(e) => setNumber(e.target.value)}
                            required
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="threshold">Threshold Percentage (optional)</Label>
                        <Input
                            id="threshold"
                            type="number"
                            step="0.01"
                            min="0"
                            max="1"
                            placeholder="e.g., 0.1 for 10%"
                            value={thresholdPercentage}
                            onChange={(e) => setThresholdPercentage(e.target.value)}
                        />
                    </div>
                    <Button type="submit" className="w-full" disabled={isLoading}>
                        {isLoading ? (
                            <>
                                <Loader2 className="mr-2 h-4 w-4 animate-spin"/>
                                Searching...
                            </>
                        ) : (
                            "Find"
                        )}
                    </Button>
                </form>
            </CardContent>
            <CardFooter className="flex flex-col items-start">
                {result && (
                    <div className="space-y-2 w-full">
                        <p>
                            <strong>Index:</strong> {result.index}
                        </p>
                        <p>
                            <strong>Value:</strong> {result.value}
                        </p>
                        <p>
                            <strong>Is Approximate:</strong> {result.is_approximate ? "Yes" : "No"}
                        </p>
                    </div>
                )}
                {error && (
                    <Alert variant="destructive" className="w-full">
                        <AlertTitle>Error</AlertTitle>
                        <AlertDescription>{error}</AlertDescription>
                    </Alert>
                )}
            </CardFooter>
        </Card>
    )
}
