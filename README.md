# Tag Github Actions

This action extract the tag version defined in the labels of the PR related to the latest commit of the current branch.

## Inputs

### `crash_on_error`
Allows to crash the action if no tag is found. Default is **true**

## Outputs

## `tag`

The tag version extracted from the labels of the PR related to the latest commit.

## Example usage

To work with this Action, you should add it to the main branch. In the Github Action's workflow, it will take the PR related to the latest commit and read the labels to find the next release version

```
  - name: Retrieve release tag from labels
    id: tag_release
    uses: quentinchampenois/gh-tag-release-action@main
    with:
      crash_on_error: "true"
  - name: Get the next version
    run: echo "This is the next version ${{ steps.tag_release.outputs.tag }}"
```

Options : 

`crash_on_error` : If set to true, the action will crash if no tag is found. Default is **true**

### Example of workflow 

You can get the tag version from the next job using the `outputs` key on job like this :

```
  ghtagrelease:
    runs-on: ubuntu-latest
    name: A job to release tags from label
    outputs:
      tag: ${{ steps.tag_release.outputs.tag }}
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Retrieve release tag from labels
        id: tag_release
        uses: quentinchampenois/gh-tag-release-action@main
        with:
          crash_on_error: "true"
      - name: Get the next version
        run: echo "This is the next version ${{ steps.tag_release.outputs.tag }}"
```

And in the next job : 

```
  tagrelease:
    runs-on: ubuntu-latest
    name: Create new tag
    needs: [ghtagrelease]
    if: ${{needs.ghtagrelease.outputs.tag}} != 'null'
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      - name: Create new tag
        id: tag
        run: |
          TAG_NAME="v${{ needs.ghtagrelease.outputs.tag }}"
          git tag -a $TAG_NAME -m "Automatic tag creation"
          echo "TAG_NAME=$TAG_NAME" >> $GITHUB_ENV
```