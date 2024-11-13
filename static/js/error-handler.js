// error-handler.js

document.addEventListener('DOMContentLoaded', function () {
   //  console.log("error-handler.js loaded"); // Debugging line to confirm script execution

    // Parse the URL for 'error' query parameter
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    console.log("Error parameter:", error); // Debugging line
    if (error) {
        // Decode the error message (in case it's URL-encoded)
        const decodedError = decodeURIComponent(error);
        console.log("Decoded Error:", decodedError); // Debugging line

        // Display the error using SweetAlert2
        Swal.fire({
            icon: 'error',
            title: 'Error',
            text: decodedError,
            confirmButtonText: 'OK',
            customClass: {
                popup: 'bg-gray-800 text-gray-200', // Tailwind classes for modal background and text
                title: 'font-semibold text-xl', // Tailwind classes for title styling
                content: 'text-gray-300', // Tailwind classes for content/body text
                confirmButton: 'bg-teal-500 hover:bg-teal-600 text-gray-900 p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-yellow-400',
                // If you have a cancel button or other buttons, style them similarly
                actions: 'flex justify-center space-x-4', // Tailwind classes for button container
                icon: 'text-yellow-400' // Optional: Customize icon color
            },
            buttonsStyling: false, // Disable default styles to apply custom Tailwind classes
        }).then(() => {
            // Remove the 'error' parameter from the URL without reloading the page
            const url = new URL(window.location);
            url.searchParams.delete('error');
            window.history.replaceState({}, document.title, url.pathname + url.search);
        });
    }
});
