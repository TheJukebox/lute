<script lang="ts">
    let files: FileList;
    function handleFiles() {
    for (let i = 0;i < files.length;i++) {
            console.log(files.item(i));
        }
    }
</script>

<!-- we need to build some sort of drop zone -->
<div class='upload-box'>
    <h1>
        Upload<br/>
        <input aria-label='upload' type="file" multiple bind:files={files}/>
        <button class="upload" aria-label="upload" on:click={handleFiles}></button>
    </h1>
</div>

<style>
	:root {
		--indigo-dye: #124e78ff;
		--goldenrod: #d5a021ff;
		--alabaster: #ede7d9ff;
		--viridian: #6a8e7fff;
		--viridian-dark: rgb(81, 121, 104);	
		--viridian-darker: rgb(61, 91, 78);	
		--bright-pink-crayola: #ea526fff;
	}

    .upload {
        width: 24px;
        aspect-ratio: 1;
        background: none;
        background-repeat: no-repeat;
        background-position: 50% 50%;
        border-radius: 50%;
        border-color: var(--goldenrod);
        border-width: 0px;
        background-image: url(./assets/upload.svg);
        filter: invert(75%);
    }

    .upload-box {
        position: relative;
        margin-bottom: 0px;
        margin-left: 25%;
        background-color: var(--alabaster);
        border-radius: 1%;
        width: 50%;
        box-shadow: 1px 1px 10px rgba(0, 0, 0, 0.1);
    }
</style>
