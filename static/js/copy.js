function copyToClipboard(element, value) {
    navigator.clipboard.writeText(value.toString()).then(() => {
        console.log("Copied to clipboard:", value);
        toastr.success('Value copied to clipboard.');

        // Add temporary feedback styling
        element.classList.add('text-green-500', 'scale-125', 'transform', 'transition', 'duration-300');

        // Remove feedback after 1 second
        setTimeout(() => {
            element.classList.remove('text-green-500', 'scale-125');
        }, 1000);
    }).catch(err => {
        toastr.error('Error copying to clipboard.');
        console.error("Failed to copy to clipboard", err);
    });
}
